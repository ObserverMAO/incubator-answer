/* eslint-disable @typescript-eslint/no-unused-vars */
/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { memo, FC } from 'react';

import classNames from 'classnames';

import Advance from '@/assets/images/mixin-advance.png';
import Prosperity from '@/assets/images/mixin-prosperity.png';
import Elite from '@/assets/images/mixin-elite.png';

interface IProps {
  /** avatar url */
  membership: string;
  /** size 48 96 128 256 */
  size: string;
  className?: string;
  alt?: string;
}

const Index: FC<IProps> = ({ membership, size, className, alt }) => {
  const roundedCls =
    className && className.indexOf('rounded') !== -1 ? '' : 'rounded';

  let membershipUrl = '';
  if (membership === 'advance') {
    membershipUrl = Advance;
  } else if (membership === 'prosperity') {
    membershipUrl = Prosperity;
  } else if (membership === 'elite') {
    membershipUrl = Elite;
  } else {
    membershipUrl = '';
  }

  return (
    <>
      {/* eslint-disable jsx-a11y/no-noninteractive-element-to-interactive-role,jsx-a11y/control-has-associated-label */}
      {membershipUrl && (
        <img
          src={membershipUrl}
          width={size}
          height={size}
          className={classNames(roundedCls, className)}
          alt={alt}
        />
      )}
    </>
  );
};

export default memo(Index);
